package com.advocatediary.tests;

import com.advocatediary.pages.DiaryPage;
import com.advocatediary.pages.SignInPage;
import org.testng.Assert;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.Test;

public class SignOutTest extends BaseTest {

  @Override
  @BeforeMethod
  public void openBrowser() {
    System.setProperty("baseUrl", "http://localhost:8080");
    System.setProperty("headless", "false");
    super.openBrowser();
  }

  @Test
  public void signOutReturnsToSignInAndBlocksDiary() {
    DiaryPage diary = signInToDiary();
    Assert.assertTrue(diary.isOpen());

    SignInPage signIn = diary.nav().signOut();
    Assert.assertTrue(signIn.isOnSignIn(), "Sign out should open the sign-in page");

    driver.get(baseUrl + "/");
    signIn.waitUntilLoaded();
    Assert.assertTrue(signIn.isOnSignIn(), "diary must stay closed after sign out");
  }
}
