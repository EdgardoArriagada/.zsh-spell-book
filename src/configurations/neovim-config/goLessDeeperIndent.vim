onoremap <silent> <bs> :<c-u>call GoLessDeeperIndent()<cr>
nnoremap <silent> <bs> :<c-u>call GoLessDeeperIndent()<cr>
vnoremap <silent> <bs> :<c-u>call GoLessDeeperIndent_Visual()<cr>

func! GoLessDeeperIndent_Visual()
  normal gv
  call GoLessDeeperIndent()
endfunc

func! GoLessDeeperIndent()
  " Go to beggin of line and add to jump list
  execute "normal!".line('.')."G^"

  :call JumpUntilNotEmptyLine('k')

  let originalInent = indent('.')

  while col('.') > 1
    let lastLineUp = GetLastMatchingIndentUp()
    execute "normal!".lastLineUp."G^"

    if indent('.') < originalInent
      break
    endif
  endwhile
endfunc

