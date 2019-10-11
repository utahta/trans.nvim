if exists('g:loaded_trans_nvim')
  finish
endif
let g:loaded_trans_nvim = 1

let s:plugin_name = 'trans'
let s:plugin_root = fnamemodify(resolve(expand('<sfile>:p')), ':h:h')
let s:plugin_cmd = [s:plugin_root . '/bin/' . s:plugin_name]

function! s:JobStart(host) abort
  return jobstart(s:plugin_cmd, { 'rpc': v:true })
endfunction

call remote#host#Register(s:plugin_name, '', function('s:JobStart'))
call remote#host#RegisterPlugin('trans', '0', [
\ {'type': 'command', 'name': 'Trans', 'sync': 1, 'opts': {'nargs': '?', 'range': '%'}},
\ {'type': 'command', 'name': 'TransWord', 'sync': 1, 'opts': {'nargs': '?', 'range': '%'}},
\ {'type': 'function', 'name': 'TransEventHandle', 'sync': 0, 'opts': {}},
\ ])
