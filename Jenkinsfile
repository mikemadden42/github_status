pipeline{
    agent any
    options {
        buildDiscarder logRotator(artifactDaysToKeepStr: '30', artifactNumToKeepStr: '2', daysToKeepStr: '30', numToKeepStr: '10')
    }
    stages{
        stage("checkout"){
            steps{
                git 'https://github.com/mikemadden42/github_status.git'
            }
        }
        stage("build"){
            steps{
                sh label: '', script: 'go build -ldflags="-s -w" -v -o github_status'
            }
        }
        stage("archive"){
            steps{
                archiveArtifacts 'github_status'
            }
        }
    }
    post{
        always{
            echo "========always========"
        }
        success{
            echo "========pipeline executed successfully ========"
        }
        failure{
            echo "========pipeline execution failed========"
        }
    }
}