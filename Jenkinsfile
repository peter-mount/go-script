properties([
  buildDiscarder(
    logRotator(
      artifactDaysToKeepStr: '',
      artifactNumToKeepStr: '',
      daysToKeepStr: '',
      numToKeepStr: '10'
    )
  ),
  disableConcurrentBuilds(),
  disableResume(),
  pipelineTriggers([
    cron('H H * * *')
  ])
])
node("go") {
  stage( 'Checkout' ) {
    checkout scm
  }
  stage( 'Init' ) {
    sh 'make clean init test'
  }
  stage( 'aix' ) {
        sh 'make -f Makefile.gen aix_ppc64'
  }
  stage( 'darwin' ) {
    parallel {
      stage( 'amd64' ) {
        sh 'make -f Makefile.gen darwin_amd64'
      }
      stage( 'arm64' ) {
        sh 'make -f Makefile.gen darwin_arm64'
      }
    }
  }
  stage( 'dragonfly' ) {
        sh 'make -f Makefile.gen dragonfly_amd64'
  }
  stage( 'freebsd' ) {
    parallel {
      stage( '386' ) {
        sh 'make -f Makefile.gen freebsd_386'
      }
      stage( 'amd64' ) {
        sh 'make -f Makefile.gen freebsd_amd64'
      }
      stage( 'arm6' ) {
        sh 'make -f Makefile.gen freebsd_arm6'
      }
      stage( 'arm64' ) {
        sh 'make -f Makefile.gen freebsd_arm64'
      }
      stage( 'arm7' ) {
        sh 'make -f Makefile.gen freebsd_arm7'
      }
    }
  }
  stage( 'illumos' ) {
        sh 'make -f Makefile.gen illumos_amd64'
  }
  stage( 'linux' ) {
    parallel {
      stage( '386' ) {
        sh 'make -f Makefile.gen linux_386'
      }
      stage( 'amd64' ) {
        sh 'make -f Makefile.gen linux_amd64'
      }
      stage( 'arm6' ) {
        sh 'make -f Makefile.gen linux_arm6'
      }
      stage( 'arm64' ) {
        sh 'make -f Makefile.gen linux_arm64'
      }
      stage( 'arm7' ) {
        sh 'make -f Makefile.gen linux_arm7'
      }
      stage( 'loong64' ) {
        sh 'make -f Makefile.gen linux_loong64'
      }
      stage( 'mips' ) {
        sh 'make -f Makefile.gen linux_mips'
      }
      stage( 'mips64' ) {
        sh 'make -f Makefile.gen linux_mips64'
      }
      stage( 'mips64le' ) {
        sh 'make -f Makefile.gen linux_mips64le'
      }
      stage( 'mipsle' ) {
        sh 'make -f Makefile.gen linux_mipsle'
      }
      stage( 'ppc64' ) {
        sh 'make -f Makefile.gen linux_ppc64'
      }
      stage( 'ppc64le' ) {
        sh 'make -f Makefile.gen linux_ppc64le'
      }
      stage( 'riscv64' ) {
        sh 'make -f Makefile.gen linux_riscv64'
      }
      stage( 's390x' ) {
        sh 'make -f Makefile.gen linux_s390x'
      }
    }
  }
  stage( 'netbsd' ) {
    parallel {
      stage( '386' ) {
        sh 'make -f Makefile.gen netbsd_386'
      }
      stage( 'amd64' ) {
        sh 'make -f Makefile.gen netbsd_amd64'
      }
      stage( 'arm6' ) {
        sh 'make -f Makefile.gen netbsd_arm6'
      }
      stage( 'arm64' ) {
        sh 'make -f Makefile.gen netbsd_arm64'
      }
      stage( 'arm7' ) {
        sh 'make -f Makefile.gen netbsd_arm7'
      }
    }
  }
  stage( 'openbsd' ) {
    parallel {
      stage( '386' ) {
        sh 'make -f Makefile.gen openbsd_386'
      }
      stage( 'amd64' ) {
        sh 'make -f Makefile.gen openbsd_amd64'
      }
      stage( 'arm6' ) {
        sh 'make -f Makefile.gen openbsd_arm6'
      }
      stage( 'arm64' ) {
        sh 'make -f Makefile.gen openbsd_arm64'
      }
      stage( 'arm7' ) {
        sh 'make -f Makefile.gen openbsd_arm7'
      }
      stage( 'mips64' ) {
        sh 'make -f Makefile.gen openbsd_mips64'
      }
    }
  }
  stage( 'plan9' ) {
    parallel {
      stage( '386' ) {
        sh 'make -f Makefile.gen plan9_386'
      }
      stage( 'amd64' ) {
        sh 'make -f Makefile.gen plan9_amd64'
      }
      stage( 'arm6' ) {
        sh 'make -f Makefile.gen plan9_arm6'
      }
      stage( 'arm7' ) {
        sh 'make -f Makefile.gen plan9_arm7'
      }
    }
  }
  stage( 'solaris' ) {
        sh 'make -f Makefile.gen solaris_amd64'
  }
  stage( 'windows' ) {
    parallel {
      stage( '386' ) {
        sh 'make -f Makefile.gen windows_386'
      }
      stage( 'amd64' ) {
        sh 'make -f Makefile.gen windows_amd64'
      }
      stage( 'arm6' ) {
        sh 'make -f Makefile.gen windows_arm6'
      }
      stage( 'arm64' ) {
        sh 'make -f Makefile.gen windows_arm64'
      }
      stage( 'arm7' ) {
        sh 'make -f Makefile.gen windows_arm7'
      }
    }
  }
}
