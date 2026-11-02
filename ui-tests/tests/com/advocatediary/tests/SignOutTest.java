package com.advocatediary.tests;

import com.advocatediary.pages.SignInPage;
import org.testng.Assert;
import org.testng.annotations.Test;

public class SignOutTest extends BaseTest {

  @Test
  public void signOutReturnsToSignInAndBlocksDiary() {
    SignInPage signIn = signOut(signInToDiary());
    Assert.assertTrue(signIn.isOnSignIn(), "Sign out should open the sign-in page");
    driver.get(baseUrl + "/");
    signIn.waitUntilLoaded();
    Assert.assertTrue(signIn.isOnSignIn(), "diary must stay closed after sign out");
  }
}
