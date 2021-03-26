{ pkgs ? import <nixpkgs> { } }:

pkgs.mkShell { PIPEDRIVE_ORG = "ruciferno-sandbox"; }
