package com.advocatediary.tests;

import com.advocatediary.pages.DiaryPage;
import org.testng.Assert;
import org.testng.annotations.Test;

public class SignupDiarySignOutTest extends BaseTest {

  @Test
  public void signupClientCaseShowsOnDiaryThenSignOut() {
    Assert.assertTrue(signupEmptyDiary().isEmpty());
    DiaryPage diary = completeClientCaseOnDiary();
    Assert.assertTrue(diaryShowsMatter(diary), "saved case should show on court diary");
    Assert.assertTrue(signOut(diary).isOnSignIn(), "Sign out should open the sign-in page");
  }
}
