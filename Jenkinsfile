node {
    def scmVars = checkout scm
    stage('Build') {
        def image = docker.build("willwangkelda/${env.JOB_NAME}:${scmVars.GIT_COMMIT}")
        image.push()
    }
    stage('Test') {
        // Skip
    }
    stage('Deploy') {
        // Skip
    }
}
