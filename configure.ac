AC_PREREQ([2.71])
AC_INIT([index], [0.0])

AC_CONFIG_AUX_DIR([build-aux])
AC_CONFIG_SRCDIR([NEWS])

AM_INIT_AUTOMAKE([subdir-objects])

GOEXT="go"
AC_SUBST([GOEXT])

AC_PROG_GO

AC_PROG_INSTALL

AC_CONFIG_FILES([src/index.go])
AC_CONFIG_FILES([Makefile])

AC_OUTPUT

cat << EOF

============================================================

${PACKAGE_NAME} Version ${PACKAGE_VERSION}

Go compiler: '${GOC} ${GOFLAGS}'

Compile and install with 'make && make install'.

============================================================

EOF
