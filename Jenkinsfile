node { 
    stage('Dependency') {
        checkout scm
        sh "make clean"
	    sh "make restore"
    }
    stage('Build') {
        checkout scm
        sh "make all"
    }
    stage('Lint') {
        checkout scm
        sh "make lint"
    }
    stage('Test') {
        checkout scm
        sh make test
    }
}
