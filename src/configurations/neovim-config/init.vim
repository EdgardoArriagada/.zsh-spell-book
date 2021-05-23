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
call Source('utils')
call Source('powerSelection')
call Source('goLessDeeperIndent')
call Source('lookDownForIndention')
call Source('plugins')

if exists('g:vscode')
  " Vscode only scripts
  call Source('vscode')
else
  " Neovim only scripts
  call Source('coc.conf')
  call Source('telescope.conf')
  call Source('theme.conf')
endif

