" Source normal vim configuration
set runtimepath^=~/.vim runtimepath+=~/.vim/after
let &packpath=&runtimepath
source ~/.vimrc

" Source helper function
let s:confPath = $HOME . '/.zsh-spell-book/src/configurations/neovim-config'
let s:source = 'source ' . s:confPath . '/'

function! Source(fileName)
   execute s:source . a:fileName . ".vim"
endfunction

" Global scripts
call Source('nvimrc')

if exists('g:vscode')
  " Vscode only scripts
  call Source('vscode')
else
  " Neovim only scripts
  call Source('plugins')
  call Source('coc.config')
endif

