pipeline {
  agent {
    docker {
      image 'golang'
    }        
  }
  stages {
    agent {
      docker {
        image 'docker'
      }
    }
    stage('Build') {
      steps {
        sh 'ls -la'
        sh 'pwd'
        sh 'export GOPATH="/go"'
        sh 'mkdir -p /go/src/github.com/atarax/rbl-control'
        sh 'cd /go/src/github.com/atarax/ && git clone https://github.com/atarax/rbl-control && cd rbl-control && go get ./...'
        sh 'go build -o bin/rbl-control'
      }
    }
    stage('Hello') {
      agent {
        docker {
          image 'docker'
        }
      }
      steps {
        sh 'echo "hello from docker"'
      }
    }
  }
}