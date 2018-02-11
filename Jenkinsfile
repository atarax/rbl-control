pipeline {
  agent {
    docker {
      image 'golang'
    }
    
  }
  stages {
    stage('Build') {
      steps {
        sh 'export GOPATH="$(pwd)/go"'
        sh 'go get ./...'
        sh 'go build -o bin/rbl-control'
      }
    }
  }
}
safd