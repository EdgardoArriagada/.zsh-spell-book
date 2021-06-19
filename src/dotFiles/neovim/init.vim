" Source normal vim configuration
set runtimepath^=~/.vim runtimepath+=~/.vim/after
let &packpath=&runtimepath
source ~/.vimrc

" Source helper function
let s:dirname = $HOME.'/.zsh-spell-book/src/dotFiles/neovim/'
let s:sourceInDir = 'source '.s:dirname

func! Source(fileName)
   exe s:sourceInDir.a:fileName.".vim"
endfunc

func! SourceFolder(folderName)
  for vimFile in split(glob(s:dirname.a:folderName.'/*.vim'), '\n')
     exe 'source '.vimFile
  endfor
endfunc

" Global scripts
call Source('nvimrc')
call Source('plugins')
call SourceFolder('custom')

if exists('g:vscode')
  call SourceFolder('vscode')
else
  call SourceFolder('neovim')
endif

