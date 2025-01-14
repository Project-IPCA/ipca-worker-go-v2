pipeline {
    agent none
    environment {
        CREDENTIALS_ID = "${env.BRANCH_NAME == 'develop' ? 'worker-dev' : 'prod'}"
        COMPOSE_FILE = "${env.BRANCH_NAME == 'develop' ? 'docker-compose.dev.yml' : 'docker-compose.yml'}"
        BUILD_OPTIONS = "${env.BRANCH_NAME == 'develop' ? '' : 'ipca-worker-go-v2 --no-deps'}"
        WORKSPACE_DIR = "${env.BRANCH_NAME == 'develop' ? '' : '/ipca/ipca-system'}"
        AGENT_NODE = "${env.BRANCH_NAME == 'develop' ? 'develop-agent' : 'master-agent'}"
    }
    options{
        skipDefaultCheckout()
    }
    stages {
        stage('Build and Deploy') {
            agent { 
                label "${AGENT_NODE}"
            }
            steps {
                script {
                    if(env.BRANCH_NAME == 'develop') {
                        checkout scm
                    }
                    withCredentials([file(credentialsId: "${CREDENTIALS_ID}", variable: 'env_file')]) {
                        if (env.BRANCH_NAME == 'develop') {
                                sh "cat ${env_file} > .env"
                                sh "docker compose -f ${COMPOSE_FILE} up -d --build"
                        } else {
                            dir("${WORKSPACE_DIR}") {
                                sh "cat ${env_file} > .env"
                                sh "cd ipca-worker-go-v2"
                                sh "git fetch"
                                sh "git pull origin main"
                                sh "cd .."
                                sh "docker compose -f ${COMPOSE_FILE} up -d --build ${BUILD_OPTIONS}"
                            }
                        }
                    }
                }
            }
        }
    }
    post {
        always {
            echo "Pipeline finished for branch: ${env.BRANCH_NAME}"
        }
    }
}