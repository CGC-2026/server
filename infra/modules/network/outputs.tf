output "vpc_id" {
  value = aws_vpc.this.id
}

output "ec2_sg_id" {
  value = aws_security_group.ec2.id
}
output "rds_sg_id" {
  value = aws_security_group.rds.id
}

output "availability_zone_a" {
  value = local.az_a
}

output "availability_zone_b" {
  value = local.az_b
}

output "private_subnet_ids" {
    value = [aws_subnet.private.id, aws_subnet.private_b.id]
}

output "public_subnet_ids" {
    value = [aws_subnet.public.id, aws_subnet.public_b.id]
}