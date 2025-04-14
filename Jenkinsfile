@Library('Jenkins-Shared-Lib') _

pipeline {
    agent {
        kubernetes {
            yaml jenkinsAgent(['agent-go': 'registry.runicrealms.com/jenkins/agent-go:latest'])
        }
    }

    environment {
        PROJECT_NAME = 'Palimpsest'
        ARTIFACT_NAME = 'palimpsest'
        REGISTRY = 'registry.runicrealms.com'
        REGISTRY_PROJECT = 'library'
    }

    stages {
        stage('Send Discord Notification (Build Start)') {
            steps {
                discordNotifyStart(env.PROJECT_NAME, env.GIT_URL, env.GIT_BRANCH, env.GIT_COMMIT.take(7))
            }
        }
        stage('Determine Environment') {
            steps {
                script {
                    def branchName = env.GIT_BRANCH.replaceAll(/^origin\//, '').replaceAll(/^refs\/heads\//, '')
                    echo "Using normalized branch name: ${branchName}"

                    if (branchName == 'dev') {
                        env.RUN_MAIN_DEPLOY = 'false'
                    } else if (branchName == 'main') {
                        env.RUN_MAIN_DEPLOY = 'true'
                    } else {
                        error "Unsupported branch: ${branchName}"
                    }
                }
            }
        }
        stage('Build and Push Server Docker Image') {
            steps {
                container('agent-go') {
                    script {
                        sh """
                        go mod download
                        go build -buildvcs=false -o palimpsest ./cmd
                        """
                        orasPush(env.ARTIFACT_NAME, "latest", "palimpsest", env.REGISTRY, env.REGISTRY_PROJECT)
                    }
                }
            }
        }
    }

    post {
        success {
            discordNotifySuccess(env.PROJECT_NAME, env.GIT_URL, env.GIT_BRANCH, env.GIT_COMMIT.take(7))
        }
        failure {
            discordNotifyFail(env.PROJECT_NAME, env.GIT_URL, env.GIT_BRANCH, env.GIT_COMMIT.take(7))
        }
    }
}
