pipeline {
  agent any

  environment {
    REGISTRY        = "registry.example.com"
    IMAGE_NAME      = "ecommerce-api"
    IMAGE_TAG       = "${env.BUILD_NUMBER}"
    DOCKER_CREDS    = credentials('dockerhub-credentials')
  }

  options {
    timestamps()
    ansiColor('xterm')
  }

  stages {
    stage('Checkout') {
      steps {
        deleteDir()
        checkout scm
      }
    }

    stage('Set Image Tag') {
      steps {
        script {
          def commit = env.GIT_COMMIT
          env.IMAGE_TAG = commit ? commit.take(7) : "build-${env.BUILD_NUMBER}"
        }
      }
    }

    stage('Unit Tests') {
      steps {
        dir('services') {
          sh 'go test ./...'
        }
      }
    }

    stage('Build Image') {
      steps {
        dir('services') {
          sh '''
            docker build -f ../docker/Dockerfile \
              -t ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} \
              -t ${REGISTRY}/${IMAGE_NAME}:latest .
          '''
        }
      }
    }

    stage('Push Image') {
      steps {
        withEnv(["DOCKER_CONFIG=${env.WORKSPACE}/.docker"]) {
          sh 'mkdir -p $DOCKER_CONFIG'
          sh "echo ${DOCKER_CREDS_PSW} | docker login ${REGISTRY} -u ${DOCKER_CREDS_USR} --password-stdin"
          sh "docker push ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
          sh "docker push ${REGISTRY}/${IMAGE_NAME}:latest"
        }
      }
    }

    stage('Deploy to Kubernetes') {
      steps {
        withCredentials([file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG_FILE')]) {
          withEnv(["KUBECONFIG=${KUBECONFIG_FILE}"]) {
            sh '''
              kubectl -n ecommerce set image deploy/ecommerce-api \
                api=${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}
              kubectl -n ecommerce rollout status deploy/ecommerce-api
            '''
          }
        }
      }
    }
  }

  post {
    always {
      sh 'docker image prune -f || true'
    }
    failure {
      mail to: 'devops@example.com',
           subject: "Jenkins build failed: ${env.JOB_NAME} ${env.BUILD_NUMBER}",
           body: "See details at ${env.BUILD_URL}"
    }
  }
}
