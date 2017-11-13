node {
    stage('Prepare') {
        checkout scm
        sh "make clean"
    }
    stage('Build') {
        checkout scm
        sh "./bin/ci-test.sh all"
    }
    stage('Lint') {
        checkout scm
        sh "./bin/ci-test.sh lint"
    }
    stage('Build and deploy') {
        checkout scm
        sh "APP=gad bash ./bin/herokutor.sh `pwd`"
    }
}
