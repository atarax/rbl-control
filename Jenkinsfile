pipeline {
  agent none
  stages {
    stage('Build Binary') {
      agent {
        docker {
          image 'golang'
        }
        
      }
      steps {
        sh "printenv"
        sh "pwd"
        sh "export GOPATH='/go'"
        sh "mkdir -p /go/src/github.com/atarax/rbl-control"
        sh '''cd /go/src/github.com/atarax/ && \
              git clone https://github.com/atarax/rbl-control && \
              cd rbl-control && \
              go get ./...'''
        sh "CGO_ENABLED=0 go build -o bin/rbl-control"
      }
    }
    stage('Build Container') {
      agent {
        docker {
          image 'docker'
        }
        
      }
      steps {
        sh "docker login -u ${DOCKERHUB_CREDENTIALS_USR} -p ${DOCKERHUB_CREDENTIALS_PSW}"
        sh "docker build . -t atarax/rbl-control:${INTERMEDIATE_IMAGE_TAG}"
        sh "docker push atarax/rbl-control:${INTERMEDIATE_IMAGE_TAG}"
      }
    }
    stage('Test Container') {
      agent {
        docker {
          image 'docker'
        }
        
      }
      steps {
        sh 'docker login -u ${DOCKERHUB_CREDENTIALS_USR} -p ${DOCKERHUB_CREDENTIALS_PSW}'
        sh '''docker run \
              -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
              -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
              atarax/rbl-control:${INTERMEDIATE_IMAGE_TAG} \
              /rbl-control -v -r "eu-west-1" -c create'''
        sh '''docker run \
              -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
              -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
              atarax/rbl-control:${INTERMEDIATE_IMAGE_TAG} \
              /rbl-control -v -r "eu-west-1" -c list'''
        sh '''docker run \
              -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
              -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
              atarax/rbl-control:${INTERMEDIATE_IMAGE_TAG} \
              /rbl-control -v -r "eu-west-1" -c destroy'''
      }
    }
    stage('Deploy Image') {
      agent {
        docker {
          image 'docker'
        }
      }
      steps {
        sh "docker login -u ${DOCKERHUB_CREDENTIALS_USR} -p ${DOCKERHUB_CREDENTIALS_PSW}"
        sh "docker pull atarax/rbl-control:${INTERMEDIATE_IMAGE_TAG}"
        sh "docker tag atarax/rbl-control:${INTERMEDIATE_IMAGE_TAG} atarax/rbl-control:${STABLE_IMAGE_TAG}"
        sh "docker push atarax/rbl-control:${STABLE_IMAGE_TAG}"
      }
    }
  }
  environment {
    DOCKERHUB_CREDENTIALS = credentials('dockerhub_credentials')
    AWS_ACCESS_KEY_ID = credentials('aws_access_key_id')
    AWS_SECRET_ACCESS_KEY = credentials('aws_secret_access_key')
    STABLE_IMAGE_TAG = "${env.BUILD_TAG + '-stable'}"
    INTERMEDIATE_IMAGE_TAG = "${env.BUILD_TAG + '-intermediate'}"
  }
}