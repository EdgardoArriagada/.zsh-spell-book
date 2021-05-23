onoremap <silent> <tab> :<c-u>call SmartTab()<cr>
nnoremap <silent> <tab> :<c-u>call SmartTab()<cr>
vnoremap <silent> <tab> :<c-u>call SmartTab_Visual()<cr>

func! SmartTab_Visual()
  normal gv
  call SmartTab()
endfunc

func! SmartTab()
  normal! ^

  " Search for a non empty to begin with
  while IsEmptyLine()
    normal! j^
  endwhile

  let l:fileBottom = line('$')
  let l:originalValidFirstCol = col('.')

  while line('.') < l:fileBottom

    normal! j^

    " [Step A] Don't use IsEmptyLine for better performance
    if len(getline(".")) == 0
      continue
    endif

    let l:currentFirstCol = col('.')

    if l:currentFirstCol > l:originalValidFirstCol
      " [Step B] Discard any false positive
      " from Step A by applying regex function this time
      if IsEmptyLine()
        continue
      endif

      return
    elseif l:currentFirstCol < l:originalValidFirstCol
      " [Step B] Discard any false positive
      " from Step A by applying regex function this time
      if IsEmptyLine()
        continue
      endif

      normal! k^
      return
    elseif l:currentFirstCol == l:originalValidFirstCol
      continue
    endif
  endwhile
endfunc
