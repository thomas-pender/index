# Configuration and Installation

Compile and install with:

1. <tt>mkdir build && cd build</tt>
2. <tt>../configure</tt>
3. <tt>make</tt>
4. <tt>sudo make install</tt>

## Configuration Options

To manulally choose a specific Go compiler and compiler flags, you simply redifine
the corresponding shell variables at configuration or compile time, e.g.,

<code><tt>../configure GOC="go" GOFLAGS="O2"</tt></code>

<tt>configure</tt> searches for a Go compiler. If none is found, then
<tt>GOC=gccgo</tt> and <tt>GOFLAGS="-g -O2"</tt> by default.

To override the default behavior manually, you may do so either at configuration
or compilation, e.g.,

<code>../configure GOC="go build" GOFLAGS="-O2"</code>

<tt>gccgo</tt> can be installed by, for example,

<code><tt>sudo apt-get install gccgo</tt></code>

# Usage

Available functions:

1. <tt>enter</tt>: enter item into index.
2. <tt>search</tt>: search for index item.
3. <tt>list</tt>: list the subjects used in index.
4. <tt>remove</tt>: remove particular item from index.

Example:
<code>
<tt>index -func=enter -author="Me MeToo" -title="Our paper"
-subject="Stuff StuffAgain"</tt>
</code>
