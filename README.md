# trans.nvim

trans.nvim is a plugin to translate text with Google Translator for Neovim.

### message
![trans_nvim_normal_low](https://user-images.githubusercontent.com/97572/35632085-05f00030-06e9-11e8-92a5-98252d71ce1a.gif)

### preview
![trans_nvim_previe_log](https://user-images.githubusercontent.com/97572/35763640-f51224d4-08f3-11e8-8d13-0510d13d240d.gif)

## Required

You must first set up authentication by creating a API Key or service account of GCP.

Then you need to install the [Go](https://golang.org/dl/) of version 1.11 or more.

### Using API Key

The API Key documentation can be found [here](https://cloud.google.com/translate/docs/auth#using_an_api_key).

Please set the environment variable `TRANS_API_KEY` to the API Key.

### Or Using Service Account

The service account documentation can be found [here](https://cloud.google.com/iam/docs/creating-managing-service-accounts).

Please set the environment variable `GOOGLE_APPLICATION_CREDENTIALS` to the file path of the JSON file that contains your service account key.

## Note

You need a little bit of money to use google translate API.  
e.g. it costs $0.06 for 2889 characters.

## Installation

For vim-plug
```viml
Plug 'utahta/trans.nvim', {'do': 'make'}
```

## Settings

```viml
let g:trans_lang_locale = 'ja'
let g:trans_lang_output = 'preview'
```

### To use floating windows

Floating windows are supported in neovim >= 0.4.0.

```viml
let g:trans_lang_output = 'float'
```

A Floating window is automatically hide when mouse cursor is moved.
