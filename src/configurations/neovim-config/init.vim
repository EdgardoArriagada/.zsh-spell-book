" Source normal vim configuration
set runtimepath^=~/.vim runtimepath+=~/.vim/after
let &packpath=&runtimepath
source ~/.vimrc

" Source helper function
let s:confPath = $HOME.'/.zsh-spell-book/src/configurations/neovim-config'
let s:sourceFromPath = 'source '.s:confPath.'/'

func! Source(fileName)
   exe s:sourceFromPath.a:fileName.".vim"
endfunc

func! SourceFolder(folderName)
  for vimFile in split(glob(s:confPath.'/'.a:folderName.'/*.vim'), '\n')
     exe 'source '.vimFile
  endfor
endfunc

" Global scripts
call Source('nvimrc')
call Source('plugins')
call SourceFolder('custom')

if exists('g:vscode')
  " Vscode only scripts
  call SourceFolder('vscode')
else
  " Neovim only scripts
  call SourceFolder('neovim')
endif

