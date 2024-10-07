{ pkgs ? import <nixpkgs> {} }:

let
  env = pkgs.mkShell {
    packages = [
      pkgs.go
      pkgs.opusTools
      pkgs.ffmpeg
      pkgs.gccgo
    ];

    env = {
      CGO_CFLAGS = "-O2 -Wno-stringop-overread";
    };
  };
in
env
