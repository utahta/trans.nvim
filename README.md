# trans.nvim

Google Translate plugin for Neovim written in Go.

## Required

You must first set up authentication by creating a service account of GCP.

The service account documentation can be found [here](https://cloud.google.com/iam/docs/creating-managing-service-accounts).

After that, set the environment variable `GOOGLE_APPLICATION_CREDENTIALS` to the file path of the JSON file that contains your service account key.

## Note

Need a little bit of money to use google translate API.

## Installation

```viml
" dein.vim
call dein#add('utahta/trans.nvim', {'build': 'make'})

" NeoBundle
NeoBundle 'utahta/trans.nvim', {'build': {'unix': 'make'}}

" vim-plug
Plug 'utahta/trans.nvim', {'do': 'make'}
```

