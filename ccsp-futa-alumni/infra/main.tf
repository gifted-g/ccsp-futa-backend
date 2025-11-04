locals {
  name = var.name
}

resource "aws_cloudwatch_log_group" "this" {
  name              = "/ecs/${local.name}"
  retention_in_days = var.log_retention_days
}

# IAM for task execution
data "aws_iam_policy_document" "task_exec_assume" {
  statement {
    actions = ["sts:AssumeRole"]
    principals { type = "Service" identifiers = ["ecs-tasks.amazonaws.com"] }
  }
}

resource "aws_iam_role" "task_exec" {
  name               = "${local.name}-exec"
  assume_role_policy = data.aws_iam_policy_document.task_exec_assume.json
}

# Attach policy for ECR pull + CW logs + SSM get parameters
resource "aws_iam_role_policy_attachment" "task_exec_ecr" {
  role       = aws_iam_role.task_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy" "task_exec_ssm" {
  name = "${local.name}-exec-ssm"
  role = aws_iam_role.task_exec.id
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect   = "Allow",
        Action   = ["ssm:GetParameters", "ssm:GetParameter", "secretsmanager:GetSecretValue"],
        Resource = "*"
      }
    ]
  })
}

# Optional task role for the app itself (e.g., S3 access)
data "aws_iam_policy_document" "task_app_assume" {
  statement {
    actions = ["sts:AssumeRole"]
    principals { type = "Service" identifiers = ["ecs-tasks.amazonaws.com"] }
  }
}

resource "aws_iam_role" "task_app" {
  name               = "${local.name}-task"
  assume_role_policy = data.aws_iam_policy_document.task_app_assume.json
}

# ECS cluster
resource "aws_ecs_cluster" "this" {
  name = local.name
}

# Security group for service
resource "aws_security_group" "svc" {
  name        = "${local.name}-sg"
  description = "Allow HTTP in, all out"
  vpc_id      = var.vpc_id

  ingress {
    description = "ALB -> ECS"
    from_port   = var.container_port
    to_port     = var.container_port
    protocol    = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# ALB + SG
resource "aws_security_group" "alb" {
  name        = "${local.name}-alb-sg"
  description = "Internet -> ALB"
  vpc_id      = var.vpc_id

  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_lb" "this" {
  name               = "${local.name}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = var.public_subnet_ids
}

resource "aws_lb_target_group" "this" {
  name        = "${local.name}-tg"
  port        = var.container_port
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = var.vpc_id

  health_check {
    path                = var.health_check_path
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 5
    matcher             = "200-399"
  }
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.this.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this.arn
  }
}

# Task definition
resource "aws_ecs_task_definition" "this" {
  family                   = local.name
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.task_cpu
  memory                   = var.task_memory
  execution_role_arn       = aws_iam_role.task_exec.arn
  task_role_arn            = aws_iam_role.task_app.arn

  container_definitions = jsonencode([
    {
      name      = local.name
      image     = var.image
      essential = true
      portMappings = [{
        containerPort = var.container_port
        hostPort      = var.container_port
        protocol      = "tcp"
      }]
      environment = [
        for k, v in var.environment : {
          name  = k
          value = v
        }
      ]
      secrets = [
        for k, arn in var.secrets_ssm_arn : {
          name      = k
          valueFrom = arn
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.this.name
          awslogs-region        = var.region
          awslogs-stream-prefix = local.name
        }
      }
    }
  ])
}

# ECS service
resource "aws_ecs_service" "this" {
  name            = local.name
  cluster         = aws_ecs_cluster.this.id
  task_definition = aws_ecs_task_definition.this.arn
  desired_count   = var.desired_count
  launch_type     = "FARGATE"

  network_configuration {
    assign_public_ip = var.assign_public_ip
    subnets          = var.private_subnet_ids
    security_groups  = [aws_security_group.svc.id]
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.this.arn
    container_name   = local.name
    container_port   = var.container_port
  }

  lifecycle {
    ignore_changes = [task_definition] # helps rolling deploys via new task defs
  }

  depends_on = [aws_lb_listener.http]
}
