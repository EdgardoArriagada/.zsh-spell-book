" Source all vim configuration
set runtimepath^=~/.vim runtimepath+=~/.vim/after
let &packpath=&runtimepath
source ~/.vimrc

" Source helper function
let s:nvimConfPath = $HOME . '/.zsh-spell-book/src/configurations/neovim-config'
function! Source(fileName)
   execute "source " . s:nvimConfPath . "/" . a:fileName . ".vim"
endfunction

if exists('g:vscode')
  call Source('vscode')
else
  call Source('nvimrc')
  call Source('plugins')
  call Source('coc.config')
endif

