set timeoutlen=200
" Map leader to which_key
nnoremap <leader> :<c-u>WhichKey '<Space>'<CR>
nnoremap <silent> <localleader> :<c-u>WhichKey '~'<CR>

" Create map to add keys to
let g:which_key_map =  {}
" Define a separator
let g:which_key_sep = '→'
" set timeoutlen=100

" Not a fan of floating windows for this
let g:which_key_use_floating_win = 0

let g:which_key_map['e'] = [':NERDTreeToggle', 'toggle view' ]

let g:which_key_map.b = {
      \ 'name' : 'buffer' ,
      \ 'r' : [':NERDTreeFind', 'reveal'],
      \ }

" Register which key map
 autocmd! User vim-which-key call which_key#register('<Space>', 'g:which_key_map')
