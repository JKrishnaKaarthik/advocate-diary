package com.advocatediary.tests;

import com.advocatediary.pages.DiaryPage;
import com.advocatediary.pages.InboxPage;
import org.testng.Assert;
import org.testng.annotations.Test;

public class SignupClientDiaryInboxTest extends BaseTest {

  @Test
  public void signupAddClientCaseThenDiaryAndInboxShowHearing() {
    Assert.assertTrue(signupEmptyDiary().isEmpty(), "new account diary starts empty");
    DiaryPage diary = completeClientCaseOnDiary();
    Assert.assertTrue(diaryShowsMatter(diary));

    InboxPage inbox = diary.nav().openInbox();
    inbox.waitForNextHearing(CLIENT_NAME, CASE_TITLE, CASE_NUMBER);
    Assert.assertTrue(inboxShowsMatter(inbox));
  }
}
