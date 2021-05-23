nnoremap <silent> _ :call GoLessDeeperIndent()<cr>
vnoremap <silent> _ :<c-u>call GoLessDeeperIndent_Visual()<cr>

func! s:strip(input)
    return substitute(a:input, '^\s*\(.\{-}\)\s*$', '\1', '')
endfunc

func! s:isEmptyLine()
    return s:strip(len(getline("."))) == 0
endfunc

func! GoLessDeeperIndent_Visual()
  normal gv
  call GoLessDeeperIndent()
endfunc

func! GoLessDeeperIndent()
  " Go to beggin of line and add to jump list
  normal! ^m'

  " Search for a non empty to begin with
  while s:isEmptyLine()
    normal! k^
  endwhile

  let l:originalValidFirstCol = col('.')

  while line('.') > 1
    let l:beforeKFirstCol = col('.')

    normal! k^

    " [Step A] Don't use s:isEmptyLine for better performance
    if len(getline(".")) == 0
      continue
    endif

    let l:currentFirstCol = col('.')

    if l:originalValidFirstCol > l:currentFirstCol && l:beforeKFirstCol > l:currentFirstCol
      " [Step B] Discard any false positive
      " from Step A by applying regex function this time
      if s:isEmptyLine()
        continue
      endif

      return
    endif
  endwhile
endfunc

