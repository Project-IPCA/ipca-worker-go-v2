pipeline {
    agent none
    stages {
        stage('Setup Environment') {
            steps {
                script {
                    echo "BRANCH_NAME: ${env.BRANCH_NAME}"
                    if (env.BRANCH_NAME == 'main') {
                        agent_label = 'master-agent'
                        docker_compose_file = 'docker-compose.prod.yml'
                    } else if (env.BRANCH_NAME == 'develop') {
                        agent_label = 'develop-agent'
                        docker_compose_file = 'docker-compose.dev.yml'
                    } else {
                        error "Branch ${env.BRANCH_NAME} is not configured!"
                    }
                }
            }
        }
        
        stage('Build and Deploy') {
            agent { label agent_label }
            steps {
                withCredentials([file(credentialsId: 'worker-dev', variable: 'env_file')]) {
                    // Set environment variables
                    sh "cat ${env_file} > .env"
                    // Start services
                    sh "docker compose -f ${docker_compose_file} up -d --build"
                }
            }
        }

        stage('Cleanup') {
            agent { label agent_label }
            options {
                skipDefaultCheckout()
            }
            steps {
                sh "docker system prune -a -f"
            }
        }
    }

    post {
        always {
            echo "Pipeline finished for branch: ${env.BRANCH_NAME}"
        }
    }
}