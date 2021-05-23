nnoremap <silent> <bs> :call LookDownForIndention()<cr>
vnoremap <silent> <bs> :<c-u>call LookDownForIndention_Visual()<cr>

func! LookDownForIndention_Visual()
  normal gv
  call LookDownForIndention()
endfunc

func! LookDownForIndention()
  " Go to beggin of line and add to jump list
  normal! ^m'

  " Search for a non empty to begin with
  while IsEmptyLine()
    normal! j^
  endwhile

  let l:originalValidFirstCol = col('.')
  let l:endOfFile = line('$')

  while line('.') < l:endOfFile

    normal! j^

    " [Step A] Don't use IsEmptyLine for better performance
    if len(getline(".")) == 0
      continue
    endif

    let l:currentFirstCol = col('.')

    if l:originalValidFirstCol == l:currentFirstCol
      " [Step B] Discard any false positive
      " from Step A by applying regex function this time
      if IsEmptyLine()
        continue
      endif

      return
    endif
  endwhile
endfunc

