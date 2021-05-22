call plug#begin('~/.config/nvim/plugged')
  Plug 'michaeljsmith/vim-indent-object'
  if !exists('g:vscode')
    Plug 'ThePrimeagen/vim-be-good'
    Plug 'fatih/vim-go', { 'tag': '*' }
    Plug 'preservim/nerdtree'
    Plug 'neoclide/coc.nvim', {'do': 'yarn install --frozen-lockfile'}
    Plug 'tpope/vim-surround'
    Plug 'airblade/vim-gitgutter'
    Plug 'tomtom/tcomment_vim'
    Plug 'mileszs/ack.vim'
    Plug 'tpope/vim-fugitive'
    Plug 'gruvbox-community/gruvbox'
    Plug 'leafgarland/typescript-vim'
    Plug 'peitalin/vim-jsx-typescript'
    Plug 'psliwka/vim-smoothie'
    Plug 'vim-airline/vim-airline'

    " telescope fuzzy finder
    Plug 'nvim-lua/popup.nvim'
    Plug 'nvim-lua/plenary.nvim'
    Plug 'nvim-telescope/telescope.nvim'
  endif
call plug#end()

