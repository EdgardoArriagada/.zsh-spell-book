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

  " Search for a non empty to begin with
  while IsEmptyLine()
    normal! k^
  endwhile

  let l:originalValidFirstCol = col('.')

  while line('.') > 1
    let l:beforeKFirstCol = col('.')

    normal! k^

    " [Step A] Don't use IsEmptyLine for better performance
    if len(getline(".")) == 0
      continue
    endif

    let l:currentFirstCol = col('.')

    if l:originalValidFirstCol > l:currentFirstCol && l:beforeKFirstCol > l:currentFirstCol
      " [Step B] Discard any false positive
      " from Step A by applying regex function this time
      if IsEmptyLine()
        continue
      endif

      return
    endif
  endwhile
endfunc

