{ mkShell, go_1_18, gopls }:
mkShell rec {
  buildInputs = [ go gopls_1_18 ];
}
