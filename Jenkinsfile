pipeline {
    agent any

    environment {
        POSTGRES_USER = 'postgres'
        POSTGRES_DB = 'ecom'
        POSTGRES_PASSWORD = credentials('postgres-password')
    }

    stages {
        stage('Checkout') {
            steps { checkout scm}
        }
        stage('Build') {
            steps { sh 'docker compose build'}
        }
        stage('Deploy') {
            steps{
                sh 'docker compose down || true'
                sh 'docker compose up -d'
            }
        }
        stage('Test'){
            steps {
                sh 'sleep 10'
                sh 'curl -f http://localhost:8000/health'
            }
        }
    }
}