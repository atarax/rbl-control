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
        sh 'ls -la'
        sh 'pwd'
        sh 'export GOPATH="/go"'
        sh 'mkdir -p /go/src/github.com/atarax/rbl-control'
        sh 'cd /go/src/github.com/atarax/ && git clone https://github.com/atarax/rbl-control && cd rbl-control && go get ./...'
        sh 'CGO_ENABLED=0 go build -o bin/rbl-control'
      }
    }
    stage('Build Container') {
      agent {
        docker {
          image 'docker'
        }
        
      }
      steps {
        sh 'docker login -u ${DOCKERHUB_CREDENTIALS_USR} -p ${DOCKERHUB_CREDENTIALS_PSW}'
        sh 'docker build . -t atarax/rbl-control:intermediate'
        sh 'docker push atarax/rbl-control:intermediate'
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
        sh 'docker run \
              -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID 
              -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY 
              atarax/rbl-control:intermediate 
              /rbl-control -r "eu-west-1" -c list'
        sh 'docker tag atarax/rbl-control:intermediate atarax/rbl-control:stable'
        sh 'docker push atarax/rbl-control:stable'
      }
    }
  }

  environment {
    DOCKERHUB_CREDENTIALS = credentials('dockerhub_credentials')
    AWS_ACCESS_KEY_ID = credentials('aws_access_key_id')
    AWS_SECRET_ACCESS_KEY = credentials('aws_secret_access_key')
  }
}