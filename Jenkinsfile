pipeline {
  agent {
    docker {
      image 'golang'
    }
    
  }
  stages {
    stage('Build') {
      steps {
        sh 'mkdir go'
        sh 'go build -o bin/rbl-control'
        sh 'export GOPATH="$(pwd)/go"'
        sh 'go get ./...'
        sh 'go get ./...'
        sh 'go build'
      }
    }
  }
}
safd