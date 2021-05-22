nnoremap <leader>ps :lua require('telescope.builtin').git_files()<cr>
nnoremap <leader>pa :lua require('telescope.builtin').find_files()<cr>
nnoremap <leader>pf :lua require('telescope.builtin').grep_string({ search = vim.fn.input("Grep For > ")})<cr>

