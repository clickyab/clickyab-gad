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
        def OUT_LOG = sh(script: 'mktemp', returnStdout: true).trim()
        def OUT_LOG_COLOR = sh(script: 'mktemp', returnStdout: true).trim()
        sh "APP=gad OUT_LOG=$OUT_LOG OUT_LOG_COLOE=$OUT_LOG_COLOR bash ./bin/herokutor.sh `pwd`"
        def color = readFile OUT_LOG_COLOR
        def message = readFile OUT_LOG
        slackSend color: color, message: message
    }
}
