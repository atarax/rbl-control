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
        sh 'go get ./...'
        sh 'go build -o bin/rbl-control'
      }
    }
  }
}
safd