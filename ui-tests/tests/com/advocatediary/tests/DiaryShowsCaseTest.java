package com.advocatediary.tests;

import com.advocatediary.pages.DiaryPage;
import org.testng.Assert;
import org.testng.annotations.Test;

public class DiaryShowsCaseTest extends BaseTest {

  @Test
  public void savedCaseDetailsShowOnCourtDiary() {
    signInToDiary();
    DiaryPage diary = completeClientCaseOnDiary();
    Assert.assertTrue(diary.isOpen());
    Assert.assertTrue(diaryShowsMatter(diary), "saved case should show on the diary");
  }
}
