if has('nvim')
  finish
endif

if exists('g:loaded_trans_nvim')
  finish
endif
let g:loaded_trans_nvim = 1

let s:trans = yarp#py3('trans_wrap')

func! Trans(...)
  return s:trans.call('trans', a:000)
endfunc

func! TransWord(...)
  return s:trans.call('trans_word', a:000)
endfunc

command! -nargs=? -range=% Trans : call Trans(<f-args>)
command! -nargs=? -range=% TransWord : call TransWord(<f-args>)

