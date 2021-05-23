call plug#begin('~/.config/nvim/plugged')
  Plug 'michaeljsmith/vim-indent-object'
  if !exists('g:vscode')
    Plug 'ThePrimeagen/vim-be-good'
    Plug 'fatih/vim-go', { 'tag': '*' }
    Plug 'neoclide/coc.nvim', {'do': 'yarn install --frozen-lockfile'}
    Plug 'mileszs/ack.vim'
    Plug 'tpope/vim-fugitive'
    Plug 'peitalin/vim-jsx-typescript'


    " Extra tooling
    Plug 'vim-airline/vim-airline'
    Plug 'airblade/vim-gitgutter'
    Plug 'preservim/nerdtree'

    " Behavior
    Plug 'psliwka/vim-smoothie'
    Plug 'tpope/vim-surround'
    Plug 'tomtom/tcomment_vim'

    " Looks
    Plug 'gruvbox-community/gruvbox'

    " Lang Support
    Plug 'rust-lang/rust.vim'
    Plug 'leafgarland/typescript-vim'

    " telescope fuzzy finder
    Plug 'nvim-lua/popup.nvim'
    Plug 'nvim-lua/plenary.nvim'
    Plug 'nvim-telescope/telescope.nvim'
  endif
call plug#end()

