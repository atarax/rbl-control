pipeline {
  agent {
    docker {
      image 'golang'
    }
    
  }
  stages {
    stage('Build') {
      steps {
        sh 'ls -la'
        sh 'pwd'
        sh 'export GOPATH="/go"'
        sh 'mkdir -p /go/src/github.com/atarax/rbl-control'
        sh 'cd /go/src/github.com/atarax/rbl-control'
        sh 'go get ./...'
        sh 'go build -o bin/rbl-control'
      }
    }
  }
}
safd