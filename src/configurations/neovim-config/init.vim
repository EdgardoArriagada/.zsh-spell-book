set runtimepath^=~/.vim runtimepath+=~/.vim/after
let &packpath=&runtimepath
source ~/.vimrc

if exists('g:vscode')
	" VSCode extension
else
	" ordinary neovim
	call plug#begin('~/.config/nvim/plugged')

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

