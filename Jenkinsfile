node {
    stage('Build and deploy') {
        checkout scm
        sh "APP=gad bash ./bin/herokutor.sh `pwd`"
    }
}
