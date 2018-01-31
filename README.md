# trans.nvim

Google Translate plugin for Neovim written in Go.

![trans_nvim_normal_low](https://user-images.githubusercontent.com/97572/35632085-05f00030-06e9-11e8-92a5-98252d71ce1a.gif)

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

## Settings

```viml
let g:trans_lang_locale = 'ja'
```

