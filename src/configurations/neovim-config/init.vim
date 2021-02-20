set runtimepath^=~/.vim runtimepath+=~/.vim/after
let &packpath=&runtimepath
source ~/.vimrc

if exists('g:vscode')
  " VSCode extension
  source $HOME/.config/nvim/vscode/settings.vim
else
  " ordinary neovim
  call plug#begin('~/.config/nvim/plugged')

  Plug 'ThePrimeagen/vim-be-good'
  Plug 'fatih/vim-go', { 'tag': '*' }
  Plug 'kien/ctrlp.vim'
  Plug 'preservim/nerdtree'
  Plug 'neoclide/coc.nvim', {'do': 'yarn install --frozen-lockfile'}
  Plug 'tpope/vim-surround'
  Plug 'jiangmiao/auto-pairs'
  Plug 'airblade/vim-gitgutter'
  Plug 'tomtom/tcomment_vim'
  Plug 'mileszs/ack.vim'
  Plug 'tpope/vim-fugitive'
  Plug 'morhetz/gruvbox'
  Plug 'leafgarland/typescript-vim'
  Plug 'peitalin/vim-jsx-typescript'

call plug#end()

endif

set clipboard=unnamedplus "Share system clipboard (require xdotool)

