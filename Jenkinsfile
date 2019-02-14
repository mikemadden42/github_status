node {

stage('checkout') {
    git 'https://github.com/mikemadden42/github_status.git'
}

stage('build') {
    sh label: '', script: 'go build -ldflags="-s -w" -v -o github_status'
}

stage('archive') {
    archiveArtifacts 'github_status'
    }
}
