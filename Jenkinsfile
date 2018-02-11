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
        sh 'ls -ls'
        sh 'docker build . -t atarax/rbl-control'
        sh 'docker push atarax/rbl-control'
      }
    }
  }
}