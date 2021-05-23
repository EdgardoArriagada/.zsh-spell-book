nnoremap <silent> _ :call GoLessDeeperIndent()<cr>
vnoremap <silent> _ :<c-u>call GoLessDeeperIndent_Visual()<cr>

func! s:strip(input)
    return substitute(a:input, '^\s*\(.\{-}\)\s*$', '\1', '')
endfunc

func! GoLessDeeperIndent_Visual()
  normal gv
  call GoLessDeeperIndent()
endfunc

func! GoLessDeeperIndent()
  " Go to beggin of line and add to jump list
  normal! ^m'
  let l:originalBeginLine = col('.')

  while line('.') > 1
    let l:beforeKLine = col('.')

    normal! k^


    if len(getline(".")) == 0
      continue
    endif

    let l:currentCol = col('.')

    if l:originalBeginLine >= l:currentCol && l:beforeKLine > l:currentCol
      " If it is NOT a blank line
      if s:strip(len(getline("."))) != 0
        return
      endif

      continue
    endif
  endwhile
endfunc

