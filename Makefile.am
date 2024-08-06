BUILT_SOURCES = $(top_builddir)/data/index.data

GOEXT=@GOEXT@

bin_PROGRAMS = src/index
nodist_src_index_SOURCES = $(top_builddir)/src/index.$(GOEXT)

$(top_builddir)/data/index.data:
	-mkdir $(top_builddir)/data
	touch $(top_builddir)/data/index.data

$(top_builddir)/src/index: $(top_builddir)/src/index.$(GOEXT) $(top_builddir)/data/index.data
	$(GOC) $(GOFLAGS) -o $@ $<

clean-local:
	-rm -rf $(top_builddir)/src

distclean-local:
	-rm -rf $(top_builddir)/data

EXTRA_DIST = $(top_srcdir)/README.md
